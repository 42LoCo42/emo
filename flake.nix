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
          mpv-unwrapped.dev
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
          vendorSha256 = "sha256-pe8Q5MHg0nJAUr+7iFNUjOB0GZClUVRGtT/ICEiuHjA=";

          inherit nativeBuildInputs buildInputs;

          prePatch = "bash generate.sh";
        };

        apps = mkApps [ "client" "server" ];

        devShell = pkgs.mkShell {
          packages = nativeBuildInputs ++ buildInputs ++ shellTools;
        };
      });
}
