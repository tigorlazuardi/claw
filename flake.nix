{
  description = "Claw development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      protoValidate =
        pkgs:
        pkgs.fetchFromGitHub {
          owner = "bufbuild";
          repo = "protovalidate";
          rev = "v1.0.0-rc.5";
          sha256 = "sha256-PTwK8+nMt7fbDrJtDj6vc/0qq8JyX1pqrtMyHnTfJ7s=";
        };
    in
    {
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          validateProtos = protoValidate pkgs;
          kittyConfig = ''
            layout tall
            launch --title="Command" fish
            launch --title="Server" wgo -file=.go clear :: go run ./cmd/claw/main.go server
            launch --title="Migrations" wgo -file=.sql clear :: go run ./cmd/goose/main.go --reset :: go run ./cmd/go-jet/main.go
            launch --hold --title="Protobuf" ${pkgs.writeShellScript "proto-watch" ''
              cd ./schemas && wgo -file=.proto -file=buf.gen.yaml -file=buf.yaml clear :: buf generate :: echo "Protobuf generated. Watching for changes..."
            ''}
          '';
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              # Go toolchain
              go_1_25

              # Database tools
              goose
              go-jet

              # Protocol Buffers
              buf
              protobuf
              protoc-gen-go
              protoc-gen-connect-go
              protoc-gen-es

              # Node.js and protobuf ES plugins
              nodejs

              # Additional useful tools
              git
              curl
              jq

              fish
              wgo

              (writeShellScriptBin "dev" ''
                ${kitty}/bin/kitty --working-directory=$KITTY_PWD --session ${pkgs.writeText "kitty.conf" kittyConfig}
              '')
            ];

            shellHook = ''
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

              export GOOSE_DBSTRING="$(pwd)/artifacts/dev.db"
              export GOROOT="${pkgs.go_1_25}/share/go"
              export KITTY_PWD="$(pwd)"
              echo "GOOSE_DBSTRING set to: $GOOSE_DBSTRING"
              echo ""

              echo "Run 'go mod tidy' to initialize dependencies"
            '';
          };
        }
      );
    };
}
