walk(if type == "object" then with_entries(select(.key | test("tenant") | not)) else . end)
