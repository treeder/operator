set -ex
# ensure working dir is clean
if [[ -z $(git status -s) ]]
then
  echo "tree is clean"
else
  echo "tree is dirty, please commit changes before running this"
  exit
fi

./build.sh

docker push treeder/operator:rsa-key-flatten
