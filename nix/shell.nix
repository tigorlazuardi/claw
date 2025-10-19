{
  perSystem =
    { pkgs, ... }:
    {
      devShells.default = pkgs.mkShell (
        let
          validateProtos = pkgs.fetchFromGitHub {
            owner = "bufbuild";
            repo = "protovalidate";
            rev = "v1.0.0-rc.5";
            sha256 = "sha256-PTwK8+nMt7fbDrJtDj6vc/0qq8JyX1pqrtMyHnTfJ7s=";
          };
        in
        {
          packages = with pkgs; [
            go_1_25

            goose
            go-jet

            buf
            protobuf
            protoc-gen-go
            protoc-gen-connect-go
            protoc-gen-es

            nodejs

            git
            curl
            jq

            fish
            wgo
            foot

            sqlitestudio

            (writeShellScriptBin "dev" ''
              systemd-run --user --scope --unit=claw-dev ${writeShellScript "dev" ''
                cd ''${PROJECT_DIR}
                foot wgo -file=.sql clear :: go run ./cmd/goose/main.go --reset :: go run ./cmd/go-jet/main.go &
                foot --working-directory=$(pwd)/webui wgo -file=svelte.config.js npm run dev &
                foot --working-directory=$(pwd)/schemas wgo -file=.proto -file=buf.gen.yaml -file=buf.yaml clear :: buf generate :: echo "Protobuf generated. Watching for changes..." &
                foot wgo -postpone -file=.go clear :: go run ./cmd/claw/main.go server &
              ''}
            '')
            (writeShellScriptBin "stop" ''
              systemctl --user stop claw-dev.scope || echo "No claw-dev scope running"
            '')
          ];
          shellHook = ''
            export PROJECT_DIR="$(git rev-parse --show-toplevel)"
            export CLAW_PROMETHEUS_ENABLE=true
            export OTEL_RESOURCE_ATTRIBUTES="service.name=http,service.namespace=claw,deployment.environment.name=local,deployment.environment=local"

            (cd "$PROJECT_DIR/webui" && npm install)

            echo "🐾 Claw development environment loaded"
            echo "Go version: $(go version)"
            echo "Node version: $(node --version)"
            echo "Available tools:"
            echo "  - go (Go compiler)"
            echo "  - goose (Database migrations)"
            echo "  - go-jet (SQL query builder)"
            echo "  - buf (Protocol buffer tool)"
            echo "  - protoc (Protocol buffer compiler)"
            echo "  - protoc-gen-go (Go protobuf plugin)"
            echo "  - protoc-gen-connect-go (ConnectRPC plugin)"
            echo "  - node/npm (JavaScript runtime and package manager)"
            echo ""

            # Setup proto validate files
            if [ ! -d "schemas/buf/validate" ]; then
              echo "Setting up proto validate files..."
              mkdir -p schemas/buf/validate
              cp -r ${validateProtos}/proto/protovalidate/buf/validate/*.proto schemas/buf/validate/
              echo "Proto validate files copied to schemas/buf/validate"
            fi

            # Create artifacts directory if it doesn't exist
            mkdir -p artifacts

            export GOOSE_DBSTRING="$(pwd)/artifacts/migrate.db"
            export GOROOT="${pkgs.go_1_25}/share/go"
            export CLAW_DATABASE__PATH="$(pwd)/artifacts/claw.db"
            export CLAW_SERVER__WEBUI__PATH="$(pwd)/cmd/claw/internal/webui"
            export CLAW_SERVER__WEBUI__DEV_MODE=true
            echo "GOOSE_DBSTRING      set to: $GOOSE_DBSTRING"
            echo "CLAW_DATABASE__PATH set to: $CLAW_DATABASE__PATH"

            echo "Run 'go mod tidy' to initialize dependencies"
          '';
        }
      );
    };
}
