clear
ng build $1 --aot
cd dist
zip -r dist.zip *.js assets/* index.html *.css
#cd ../
#cp dist.zip ../../out/dist.zip
