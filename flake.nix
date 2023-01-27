{
  description = "Composable command line toolkit";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-22.11";
  };

  outputs = { self, nixpkgs }:
    let
      # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      # System types to support.
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    rec {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          kit = pkgs.buildGoModule {
            pname = "kit";
            inherit version;

            src = pkgs.lib.cleanSourceWith {
              src = ./.;
              filter = path: type:
              let
                p = baseNameOf path;
              in !(
                p == "flake.nix" ||
                p == "flake.lock" ||
                p == "README.md"
              );
            };

            vendorSha256 = "sha256-hjvL4WI+QzWWxXYj2tKwsYmBEkQ6V37YbIWwn/92xfI=";
          };
        });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.stdenv.mkDerivation {
            name = "kit";
            buildInputs = [
              pkgs.gotools
            ] ++ packages.${system}.kit.nativeBuildInputs;
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.kit);
    };
}
