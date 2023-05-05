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
        packages = with pkgs; [
          bashInteractive
          libsodium.dev
          oapi-codegen
          ref-merge.outputs.defaultPackage.${system}
          yq
        ];
      in
      {
        defaultPackage = pkgs.buildGoModule {
          pname = "emo";
          version = "1.0.0-indev";
          src = ./.;

          vendorSha256 = pkgs.lib.fakeSha256;

          nativeBuildInputs = packages;

          PKG_CONFIG_PATH = "${pkgs.libsodium.dev}/lib/pkgconfig";
        };

        apps = mkApps [ "client" "server" ];

        devShell = pkgs.mkShell {
          packages = packages ++ (with pkgs; [
            go
            gopls
          ]);
        };
      });
}
