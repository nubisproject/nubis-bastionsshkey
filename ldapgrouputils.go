package main

func UserInGroup(username string, group []LDAPUserObject) bool {
	for _, g_entry := range group {
		if g_entry.Uid == username {
			return true
		}
	}
	return false
}

func GetLDAPUserObjectFromGroup(username string, group []LDAPUserObject) (LDAPUserObject, bool) {
	for _, g_entry := range group {
		if g_entry.Uid == username {
			return g_entry, true
		}
	}
	return LDAPUserObject{}, false
}

func IgnoreUserLDAPUserObjects(s []LDAPUserObject, e string) bool {
	for _, a := range s {
		if a.Uid == e {
			return true
		}
	}
	return false
}

func IgnoreUser(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
