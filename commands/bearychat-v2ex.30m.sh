#!/bin/bash

curl "https://www.googleapis.com/customsearch/v1?key=AIzaSyC1Q3F9GsEaIaxLe4zRwMeOhhNr7axtXEg&cx=011777316675351136864:22g5hinnt0i&q=bearychat" | jq '.items [] | .title, .link'