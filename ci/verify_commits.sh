git rebase origin/main --exec "make check"
status=$?
if [ $status -ne 0 ]; then
  echo ""
  echo "Failing commit:"
  echo ""
  git show --name-only
  exit 1
fi

exit 0
