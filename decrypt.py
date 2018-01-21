import pyAesCrypt
# encryption/decryption buffer size - 64K
bufferSize = 64 * 1024
password = password = input("Enter password: ")
# decrypt
pyAesCrypt.decryptFile("hass-go-calendar.TOML.aes", "hass-go-calendar.TOML", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-sensors.TOML.aes", "hass-go-sensors.TOML", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-slack.TOML.aes", "hass-go-slack.TOML", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-weather.TOML.aes", "hass-go-weather.TOML", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-secrets.TOML.aes", "hass-go-secrets.TOML", password, bufferSize)
