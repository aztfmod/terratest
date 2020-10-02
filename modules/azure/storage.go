package azure

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/require"
)

// StorageAccountExists indicates whether the storage account name exactly matches; otherwise false.
// This function would fail the test if there is an error.
func StorageAccountExists(t *testing.T, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	result, err := StorageAccountExistsE(storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// StorageBlobContainerExists returns true if the container name exactly matches; otherwise false
// This function would fail the test if there is an error.
func StorageBlobContainerExists(t *testing.T, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	result, err := StorageBlobContainerExistsE(containerName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// GetStorageBlobContainerPublicAccess indicates whether a storage container has public access; otherwise false.
// This function would fail the test if there is an error.
func GetStorageBlobContainerPublicAccess(t *testing.T, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	result, err := GetStorageBlobContainerPublicAccessE(containerName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// GetStorageAccountKind returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
// This function would fail the test if there is an error.
func GetStorageAccountKind(t *testing.T, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	result, err := GetStorageAccountKindE(storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// GetStorageAccountSkuTier returns the storage account sku tier as Standard or Premium.
// This function would fail the test if there is an error.
func GetStorageAccountSkuTier(t *testing.T, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	result, err := GetStorageAccountSkuTierE(storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// GetStorageDNSString builds and returns the storage account dns string if the storage account exists.
// This function would fail the test if there is an error.
func GetStorageDNSString(t *testing.T, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	result, err := GetStorageDNSStringE(storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// StorageAccountExistsE indicates whether the storage account name exactly matches; otherwise false.
func StorageAccountExistsE(storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return false, err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return false, err2
	}
	storageAccount, err3 := GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID)
	if err3 != nil {
		return false, nil
	}
	return *storageAccount.Name == storageAccountName, nil
}

// StorageBlobContainerExistsE returns true if the container name exactly matches; otherwise false.
func StorageBlobContainerExistsE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return false, err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return false, err2
	}
	container, err := GetStorageBlobContainerE(containerName, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return (*container.Name == containerName), nil
}

// GetStorageBlobContainerPublicAccessE indicates whether a storage container has public access; otherwise false.
func GetStorageBlobContainerPublicAccessE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return false, err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return false, err2
	}
	client, err := GetStorageBlobContainerClientE(subscriptionID)
	if err != nil {
		return false, err
	}
	container, err := client.Get(context.Background(), resourceGroupName, storageAccountName, containerName)
	if err != nil {
		return false, err
	}
	return (string(container.PublicAccess) != "None"), nil
}

// GetStorageAccountKindE returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
func GetStorageAccountKindE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return "", err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return "", err2
	}
	storageAccount, err := GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}
	return string(storageAccount.Kind), nil
}

// GetStorageAccountSkuTierE returns the storage account sku tier as Standard or Premium.
func GetStorageAccountSkuTierE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}
	return string(storageAccount.Sku.Tier), nil
}

// GetStorageBlobContainerE returns Blob container client.
func GetStorageBlobContainerE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (*storage.BlobContainer, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return nil, err2
	}
	client, err := GetStorageBlobContainerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	container, err := client.Get(context.Background(), resourceGroupName, storageAccountName, containerName)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

// GetStorageAccountPropertyE returns StorageAccount properties.
func GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID string) (*storage.Account, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return nil, err2
	}
	client, err := GetStorageAccountClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	account, err := client.GetProperties(context.Background(), resourceGroupName, storageAccountName, "")
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetStorageAccountClientE creates a storage account client.
func GetStorageAccountClientE(subscriptionID string) (*storage.AccountsClient, error) {
	storageAccountClient := storage.NewAccountsClient(os.Getenv(AzureSubscriptionID))
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	storageAccountClient.Authorizer = *authorizer
	return &storageAccountClient, nil
}

// GetStorageBlobContainerClientE creates a storage container client.
func GetStorageBlobContainerClientE(subscriptionID string) (*storage.BlobContainersClient, error) {
	blobContainerClient := storage.NewBlobContainersClient(os.Getenv(AzureSubscriptionID))
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}
	blobContainerClient.Authorizer = *authorizer
	return &blobContainerClient, nil
}

// GetStorageURISuffixE returns the proper storage URI suffix for the configured Azure environment.
func GetStorageURISuffixE() (*string, error) {
	var envName = "AzurePublicCloud"
	env, err := azure.EnvironmentFromName(envName)
	if err != nil {
		return nil, err
	}
	return &env.StorageEndpointSuffix, nil
}

// GetStorageAccountPrimaryBlobEndpointE gets the storage account blob endpoint as URI string.
func GetStorageAccountPrimaryBlobEndpointE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *storageAccount.AccountProperties.PrimaryEndpoints.Blob, nil
}

// GetStorageDNSStringE builds and returns the storage account dns string if the storage account exists.
func GetStorageDNSStringE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	retval, err := StorageAccountExistsE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}
	if retval {
		storageSuffix, _ := GetStorageURISuffixE()
		return fmt.Sprintf("https://%s.blob.%s/", storageAccountName, *storageSuffix), nil
	}

	return "", NewNotFoundError("storage account", storageAccountName, "")
}