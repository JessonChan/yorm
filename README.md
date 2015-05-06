# README #

yOrm is a simple,lightweight orm  , for mysql only now.

### Why this project calls yOrm ###

yOrm is just a name.
more about the detail,cc [https://github.com/lewgun]

### What is this yOrm for? ###

* A simple mysql orm to crud

## Tags ##
 
Now support these types of tag.
### column ###
this tag alias struct name to a real column name. "Id int \`yOrm:column(autoId)\`" means this field Id will name autoId in mysql column
### pk ###
this tag allow you to set a primary key where select/delete/update as the where clause  "Id int \`yOrm:column(autoId);pk\`"




