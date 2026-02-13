.data.admin.tenants.nodes |= (map(if .name == "TestTenant" then .accountId = "test_account_id" else . end) | sort_by(.name))
