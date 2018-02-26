clear 
ng build --env=dev --aot
cd dist
rm *.map
zip ../dist.zip *.js assets/* index.html
cd ../
