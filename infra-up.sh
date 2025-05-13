
az deployment group create \
    --resource-group rg-acc-laptime-tracker \
    --template-file ./main.bicep \
    --parameters \
        cosmosDbConnectionString="$(pass cosmos/acclaptracker)" \
        githubToken="$(pass acc_github_pat)"
