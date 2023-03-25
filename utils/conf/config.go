package conf

import "github.com/spf13/viper"

// Config menyimpan semua konfigurasi dari aplikasi
// nilai nya berasal dari viper yang membaca dari config file atau environtment variable
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig membaca configuration dari file atau environtment variables
func LoadConfig(path string) (config Config, err error) {
	// SETUP VIPER UNTUK MEMBACA DARI FILE
	viper.AddConfigPath(path) // memberitahu viper lokasi configuration file nya
	viper.SetConfigName("app") // memberitahu viper untuk mencari file dengan nama yang  ada di argument "app" dari app.env
	viper.SetConfigType("env") // memberitahu tipe filenya

	// SETUP VIPER UNTUK MEMBACA DARI ENVIRONTMENT. JIKA ADA MAKA OVERRIDE SETUP DIATAS
	viper.AutomaticEnv() // otomatis mengoverride nilai dari file configuration jika ada configuration dari environment
	
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config) // mengkonversi config dari file atau env ke target 
	return
}