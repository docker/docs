.PHONY:
check:
	for img in $$(git ls-files images); do	\
	  base=$$(basename "$$img");		\
	  if ! git grep -q -F "$$base"; then	\
	    echo >&2 "unused image: $$img";	\
	    pass=false;				\
	  fi;					\
	done;					\
	$$pass
