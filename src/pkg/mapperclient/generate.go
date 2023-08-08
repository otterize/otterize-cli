package mapperclient

//go:generate wget https://raw.githubusercontent.com/otterize/network-mapper/main/src/mappergraphql/schema.graphql -O ./schema.graphql -q
//go:generate go run github.com/Khan/genqlient@v0.5.0 ./genqlient.yaml
