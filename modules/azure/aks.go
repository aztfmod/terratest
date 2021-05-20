package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-11-01/containerservice"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetManagedClustersClientE is a helper function that will setup an Azure ManagedClusters client on your behalf
func GetManagedClustersClientE(subscriptionID string) (*containerservice.ManagedClustersClient, error) {
	// Create a cluster client
	client, err := CreateManagedClustersClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// setup authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer
	return &client, nil
}

// GetManagedClusterE will return ManagedCluster
func GetManagedClusterE(t testing.TestingT, resourceGroupName, clusterName, subscriptionID string) (*containerservice.ManagedCluster, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	client, err := GetManagedClustersClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	managedCluster, err := client.Get(context.Background(), resourceGroupName, clusterName)
	if err != nil {
		return nil, err
	}
	return &managedCluster, nil
}

// ManagedClusterVersionMatch indicates whether the expectued Kuberntes Version is deployed.
// This function returns false and a non nil error object if an error occurs.
func ManagedClusterVersionMatch(t testing.TestingT, k8sVersion, resourceGroupName, clusterName, subscriptionID string) (bool, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return false, err
	}
	client, err := GetManagedClustersClientE(subscriptionID)
	if err != nil {
		return false, err
	}
	managedCluster, err := client.Get(context.Background(), resourceGroupName, clusterName)
	if err != nil {
		return false, err
	}
	return (*managedCluster.ManagedClusterProperties.KubernetesVersion == k8sVersion), nil
}
