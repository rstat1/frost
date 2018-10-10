clear

ng build $1 --aot
cd dist
rm *.map
zip ../dist.zip *.js assets/* index.html *.css
cd ../
