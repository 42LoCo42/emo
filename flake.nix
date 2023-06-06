{
  description = "easy music organizer";

  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.ref-merge.url = "github:42loco42/ref-merge";
  inputs.ref-merge.inputs.nixpkgs.follows = "nixpkgs";
  inputs.ref-merge.inputs.flake-utils.follows = "flake-utils";

  outputs = { self, nixpkgs, flake-utils, ref-merge }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        mkApps = (names: builtins.listToAttrs (map
          (name: {
            inherit name;
            value = flake-utils.lib.mkApp {
              drv = self.defaultPackage.${system};
              exePath = "/bin/${name}";
            };
          })
          names));

        nativeBuildInputs = with pkgs; [
          oapi-codegen
          pkg-config
          ref-merge.outputs.defaultPackage.${system}
          yq
        ];

        buildInputs = with pkgs; [
          libsodium
          mpv-unwrapped.dev
        ];

        shellTools = with pkgs; [
          bashInteractive
          go
          gopls
          tree
        ];
      in
      {
        defaultPackage = pkgs.buildGoModule {
          pname = "emo";
          version = "1.0.0-indev";
          src = ./.;
          vendorSha256 = "sha256-LP0VPGzvkHAhUyfsf9JnyRUHYYmHBLKaAGTH3u9Ekk4=";

          inherit nativeBuildInputs buildInputs;

          prePatch = "bash generate.sh";
        };

        apps = mkApps [ "cli" "daemon" "server" ];

        devShell = pkgs.mkShell {
          packages = nativeBuildInputs ++ buildInputs ++ shellTools;
        };
      });
}
