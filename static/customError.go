package static

type error interface {
	Error() string
}

/*
	Error pada saat user tidak ditemukan.
	Digunakan pada saat login.
*/
type AuthError struct {
}

func (e *AuthError) Error() string {
	return "User tidak ditemukan"
}

/*
	Error pada saat user tidak terautorisasi.
	Digunakan pada saat verif jwt.
*/
type Unauthorized struct {
}

func (e *Unauthorized) Error() string {
	return "Anda tidak berhak mengakses data ini."
}
