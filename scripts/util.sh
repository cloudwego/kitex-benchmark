function check_supported_env() {
  case "$OSTYPE" in
      linux*)   ;;
      darwin*)  ;;
      *)        echo "[ERROR] kitex benchmark is not supported on $OSTYPE"; exit 1;;
  esac
}