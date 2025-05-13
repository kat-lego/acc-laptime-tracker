@description('Location for all resources')
param location string = resourceGroup().location

@secure()
@description('Cosmos DB connection string')
param cosmosDbConnectionString string

@description('Cosmos DB database name')
param cosmosDbDatabase string = 'sessions'

@description('Cosmos DB container name')
param cosmosDbContainer string = 'sessions'

@secure()
@description('github pat token')
param githubToken string

resource appServicePlan 'Microsoft.Web/serverfarms@2024-04-01' = {
  name: 'asp-acc-laptime-tracker'
  location: location
  sku: {
    name: 'F1'
    capacity: 1
  }
  kind: 'linux'
  properties: {
    reserved: true
  }
}

resource apiWebApp 'Microsoft.Web/sites@2024-04-01' = {
  name: 'app-acc-laptime-tracker-api'
  location: location
  kind: 'app'
  properties: {
    serverFarmId: appServicePlan.id
    siteConfig: {
      linuxFxVersion: 'DOCKER|99percentmoses/acc-laptime-tracker-api:latest'
      appSettings: [
        {
          name: 'WEBSITES_PORT'
          value: '80'
        }
        {
          name: 'PORT'
          value: '80'
        }
        {
          name: 'ACC_COSMOS_CONNECTION_STRING'
          value: cosmosDbConnectionString
        }
        {
          name: 'ACC_COSMOS_DATABASE'
          value: cosmosDbDatabase
        }
        {
          name: 'ACC_COSMOS_CONTAINER'
          value: cosmosDbContainer
        }
        {
          name: 'ACC_CORS_ORIGINS'
          value: 'https://acc.api.katlegomodupi.com'
        }
      ]
    }
    httpsOnly: true
  }
  identity: {
    type: 'SystemAssigned'
  }
}

resource staticWebApp 'Microsoft.Web/staticSites@2024-04-01' = {
  name: 'swa-acc-laptime-tracker'
  location: 'westeurope'
  sku: {
    name: 'Free'
    tier: 'Free'
  }
  properties: {
    repositoryUrl: 'https://github.com/kat-lego/acc-laptime-tracker'
    branch: 'main'
    repositoryToken: githubToken
    buildProperties: {
      appLocation: 'web'
      appBuildCommand: 'npm run build'
      outputLocation: 'out'
    }
    provider: 'GitHub'
  }
}
