{
  description = "easy music organizer";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
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
          pkg-config
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
          vendorSha256 = "sha256-srsgCS/2dYR6WesjMQ+dgtwBc2FyIzyQWXTaiDVM3sg=";

          inherit nativeBuildInputs buildInputs;

          # prePatch = "bash generate.sh";
        };

        apps = mkApps [ "cli" "daemon" "server" ];

        devShell = pkgs.mkShell {
          packages = nativeBuildInputs ++ buildInputs ++ shellTools;
        };
      });
}
