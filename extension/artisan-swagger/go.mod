module github.com/Myriad-Dreamin/artisan/extension/artisan-swagger

go 1.14

replace github.com/Myriad-Dreamin/artisan/artisan-core v0.0.0-20200929180408-b26f5ec688a9 => ../../artisan-core

require (
	github.com/Myriad-Dreamin/artisan/artisan-core v0.0.0-20200929180408-b26f5ec688a9
	github.com/Myriad-Dreamin/minimum-lib v0.0.0-20200719050009-6377966ced3b
	github.com/go-openapi/spec v0.19.9
	github.com/stretchr/testify v1.4.0
)
