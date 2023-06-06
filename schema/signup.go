package schema

type SignUp struct {
	Nik string `json:"nik" validate:"required"`
	// Username     string `json:"username" validate:"required"`
	Nama         string `json:"nama" validate:"required"`
	Alamat       string `json:"alamat" validate:"required"`
	JenisKelamin string `json:"jenis_kelamin" validate:"required"`
	// TanggalLahir string `json:"tanggal_lahir" validate:"required"`
	NoWa      string `json:"no_wa" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Role      string `json:"role" validate:"required"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
