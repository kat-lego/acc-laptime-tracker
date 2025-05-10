
az deployment group create \
    --resource-group rg-acc-laptime-tracker \
    --template-file ./main.bicep \
    --parameters \
        tuiSshHostKeyPEM="$(cat ~/.ssh/id_acc_laptime_tracker_tui)" \
        cosmosDbConnectionString="$(pass cosmos/acclaptracker)"
