# README #

yorm is a simple,lightway orm lib , for mysql only now.

### What is this yorm for? ###

* A simple mysql orm to crud

## Tags ##
 
Now support these types of tag.
### column ###
this tag is to alias struct name to a real column name. "Id int \`yorm:column(autoId)\`" means this field Id will name autoId in mysql column
### pk ###
this tag allow you to set a primary key where select/delete/update as the where clause  "Id int \`yorm:column(autoId);pk\`"




