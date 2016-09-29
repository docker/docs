provider hyperkit {
	probe vmx__exit(int, unsigned int);
	probe vmx__ept__fault(int, unsigned long, unsigned long);
	probe vmx__inject__virq(int, int);
	probe vmx__write__msr(int, unsigned int, unsigned long);
	probe vmx__read__msr(int, unsigned int, unsigned long);

	probe block__preadv(off_t, size_t);
	probe block__preadv__done(off_t, ssize_t);
	probe block__pwritev(off_t, size_t);
	probe block__pwritev__done(off_t, ssize_t);
};
