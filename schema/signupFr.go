package schema

type SignUpFreelance struct {
	Nik          string  `json:"nik" validate:"required"`
	Nama         string  `json:"nama" validate:"required"`
	JenisKelamin string  `json:"jenis_kelamin" validate:"required"`
	Dob          string  `json:"dob" validate:"required"`
	NoWa         string  `json:"no_wa" validate:"required"`
	Email        string  `json:"email" validate:"required"`
	Password     string  `json:"password" validate:"required"`
	Role         string  `json:"role" validate:"required"`
	AddressLong  float64 `json:"address_long" validate:"required"`
	AddressLat   float64 `json:"address_lat" validate:"required"`
	Address      string  `json:"address" validate:"required"`
	JobChildCode string  `json:"kategori_pekerjaan" validate:"required"`
}

/*
{
    "nik": "09876",
    "nama": "dan kuroto",
    "jenis_kelamin": "true",
    "dob": "2001-09-09",
    "no_wa": "0987654321",
    "email": "dankuroto@gmail.com",
    "password": "password",
    "role": "freelancer",
    "address": "test",
    "address_long": 1.0,
    "address_lat": 0.1,
    "kategori_pekerjaan": "SE-TV"
}
*/
