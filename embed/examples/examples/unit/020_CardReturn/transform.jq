walk(if type == "object" then with_entries(select(.key | test("created") or test("modified") | not)) else . end)
