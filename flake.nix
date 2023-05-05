{
  description = "easy music organizer";

  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.ref-merge.url = "github:42loco42/ref-merge";

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
        ];

        shellTools = with pkgs; [
          bashInteractive
          go
          gopls
        ];
      in
      {
        defaultPackage = pkgs.buildGoModule {
          pname = "emo";
          version = "1.0.0-indev";
          src = ./.;
          vendorSha256 = "sha256-5lp8LA9KpyFqlA15gg3+/61BT8C1V9HJQw0dh/u3IRk=";

          inherit nativeBuildInputs buildInputs;

          prePatch = "bash generate.sh";
        };

        apps = mkApps [ "client" "server" ];

        devShell = pkgs.mkShell {
          packages = nativeBuildInputs ++ buildInputs ++ shellTools;
        };
      });
}
