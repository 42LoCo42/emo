{
  description = "easy music organizer";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
  flake-utils.lib.eachDefaultSystem (system: let
    pkgs   = import nixpkgs { inherit system; };
    mkApps = (names: builtins.listToAttrs (map (name: {
      inherit name;
      value = flake-utils.lib.mkApp {
        drv = self.defaultPackage.${system};
        exePath = "/bin/${name}";
      };
    }) names));
  in {
    defaultPackage = pkgs.buildGoModule {
      pname = "emo";
      version = "0.0.1";
      src = ./.;

      vendorSha256 = "Qf1JPyqBsqlcCQ1/hjYOEOsbsWodm6ETZSv6F/Aj9DQ=";

      nativeBuildInputs = with pkgs; [
        pkg-config
        libsodium.dev
      ];

      PKG_CONFIG_PATH = "${pkgs.libsodium.dev}/lib/pkgconfig";
    };

    apps = mkApps ["client" "server"];
  });
}
