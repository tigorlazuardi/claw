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
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              # Go toolchain
              go

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

              export GOOSE_DBSTRING="artifacts/dev.db"
              echo "GOOSE_DBSTRING set to: $GOOSE_DBSTRING"
              echo ""

              echo "Run 'go mod tidy' to initialize dependencies"
            '';
          };
        }
      );
    };
}
