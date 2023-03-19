#!/bin/bash

set -euo pipefail

ls $LEXICON_PATH

for lex in "NWL20" "NWL18" "America" "CSW21" "CSW19" "ECWL"
do
    awk '{print $1}' "$LEXICON_PATH/$lex.txt" > "$LEXICON_PATH/$lex-stripped.txt"
    awk '{print toupper($0)}' "$LEXICON_PATH/$lex-stripped.txt" > "$LEXICON_PATH/$lex-toupper.txt"
    echo "lex $lex"

    CONTAINER_ID="$(docker create kbuilder -- english-kwg /home/in.txt /home/out.kwg )"
    trap "docker rm $CONTAINER_ID" EXIT
    echo "$CONTAINER_ID"

    docker cp "$LEXICON_PATH/$lex-toupper.txt" "$CONTAINER_ID:/home/in.txt"
    docker start "$CONTAINER_ID"
    docker attach "$CONTAINER_ID" || true
    docker cp "$CONTAINER_ID:/home/out.kwg" "$LEXICON_PATH/gaddag/$lex.kwg"

    docker rm "$CONTAINER_ID"
    trap "" EXIT

    echo "after $lex"
done

for lex in "OSPS44"
do
    awk '{print $1}' "$LEXICON_PATH/$lex.txt" > "$LEXICON_PATH/$lex-stripped.txt"
    echo "lex $lex"

    CONTAINER_ID="$(docker create kbuilder -- polish-kwg /home/in.txt /home/out.kwg )"
    trap "docker rm $CONTAINER_ID" EXIT
    echo "$CONTAINER_ID"

    docker cp "$LEXICON_PATH/$lex-stripped.txt" "$CONTAINER_ID:/home/in.txt"
    docker start "$CONTAINER_ID"
    docker attach "$CONTAINER_ID" || true
    docker cp "$CONTAINER_ID:/home/out.kwg" "$LEXICON_PATH/gaddag/$lex.kwg"

    docker rm "$CONTAINER_ID"
    trap "" EXIT

    echo "after $lex"
done


for lex in "NSF22"
do
    awk '{print $1}' "$LEXICON_PATH/$lex.txt" > "$LEXICON_PATH/$lex-stripped.txt"
    awk '{print toupper($0)}' "$LEXICON_PATH/$lex-stripped.txt" > "$LEXICON_PATH/$lex-toupper.txt"
    echo "lex $lex"

    CONTAINER_ID="$(docker create kbuilder -- norwegian-kwg /home/in.txt /home/out.kwg )"
    trap "docker rm $CONTAINER_ID" EXIT
    echo "$CONTAINER_ID"

    docker cp "$LEXICON_PATH/$lex-toupper.txt" "$CONTAINER_ID:/home/in.txt"
    docker start "$CONTAINER_ID"
    docker attach "$CONTAINER_ID" || true
    docker cp "$CONTAINER_ID:/home/out.kwg" "$LEXICON_PATH/gaddag/$lex.kwg"

    docker rm "$CONTAINER_ID"
    trap "" EXIT

    echo "after $lex"
done

echo "done creating kwgs"
