walk(if type == "object" then with_entries(select(.key | test("fileId") | not)) else . end)
