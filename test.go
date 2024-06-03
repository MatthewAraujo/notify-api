curl -L \
  -X DELETE \
  eH "Accept: application/vnd.github+json" \
-H "Authorization: Bearer $token" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/MatthewAraujo/test/hooks/482148661
