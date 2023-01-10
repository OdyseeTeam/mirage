module github.com/OdyseeTeam/mirage

go 1.18

require (
	github.com/Depado/ginprom v1.7.4
	github.com/OdyseeTeam/gody-cdn v1.0.6
	github.com/bluele/gcache v0.0.2
	github.com/chai2010/webp v1.1.1
	github.com/ekyoung/gin-nice-recovery v0.0.0-20160510022553-1654dca486db
	github.com/gabriel-vasile/mimetype v1.4.1
	github.com/gin-contrib/pprof v1.4.0
	github.com/gin-gonic/gin v1.8.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/h2non/bimg v1.1.9
	github.com/johntdyer/slackrus v0.0.0-20220912135606-861993969176
	github.com/lbryio/lbry.go/v2 v2.7.2-0.20220815204100-2adb8af5b68c
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/oov/psd v0.0.0-20220121172623-5db5eafcecbb
	github.com/prometheus/client_golang v1.13.0
	github.com/sirupsen/logrus v1.9.0
	github.com/sizeofint/gif-to-webp v0.0.0-20210224202734-e9d7ed071591
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.13.0
	golang.org/x/image v0.0.0-20220902085622-e7cb96979f69
)

require (
	github.com/aws/aws-sdk-go v1.44.104 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/brk0v/directio v0.0.0-20190225130936-69406e757cf7 // indirect
	github.com/c2h5oh/datasize v0.0.0-20200825124411-48ed595a09d2 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.1 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190915194858-d3ddacdb130f // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/johntdyer/slack-go v0.0.0-20180213144715-95fac1160b22 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/karrick/godirwalk v1.17.0 // indirect
	github.com/lbryio/reflector.go v1.1.3-0.20220730181028-f5d30b1a6e79 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/sizeofint/webp-animation v0.0.0-20190207194838-b631dc900de9 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/tkanos/gonfig v0.0.0-20210106201359-53e13348de2f // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sync v0.0.0-20220907140024-f12130a52804 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/nullbio/null.v6 v6.0.0-20161116030900-40264a2e6b79 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/OdyseeTeam/gody-cdn => /home/niko/work/repositories/gody-cdn/
replace github.com/lbryio/lbry.go/v2 => github.com/OdyseeTeam/lbry.go/v2 v2.7.2-0.20221101212832-e6a3f4002995
