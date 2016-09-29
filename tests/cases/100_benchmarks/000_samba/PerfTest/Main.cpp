#include <chrono>
#include <vector>
#include <thread>
#include <string>
#include <iostream>
#include <memory>
#include <cstdint>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <cstring>

using namespace std;
void worker(const string& baseDir, int threadIndex, int fileCount, int blockSize, uint8_t* buf, int blocksPerFile);
int main(int argc, const char* argv[]) {
	if (argc != 6) {
		cerr << "usage : PerfTest <threads> <files_per_thread> <block_size> <blocks_per_file> <base_directory>" << endl;
		return -1;
	}
	// parsing args
	auto threadCount = std::stoi(argv[1]);
	auto filesPerThread = std::stoi(argv[2]);
	auto blockSize = std::stoi(argv[3]);
	auto blocksPerFile = std::stoi(argv[4]);
	string baseDirectory = argv[5];

	if (threadCount < 1) {
		cerr << "threads must be >= 1" << endl;
		return -1;
	}
	if (filesPerThread < 1) {
		cerr << "files_per_thread must be >= 1" << endl;
		return -1;
	}
	if (blockSize < 1) {
		cerr << "block_size must be >= 1" << endl;
		return -1;
	}
	if (blocksPerFile < 1) {
		cerr << "blocks_per_file must be >= 1" << endl;
		return -1;
	}
	if (baseDirectory.size() == 0 || baseDirectory[0] != '/') {
		cerr << "base_directory must be an absolute directory" << endl;
		return -1;
	}

	// preparing buffer
	auto buffer = unique_ptr<uint8_t[]>(new uint8_t[blockSize]);
	for (auto ix = 0; ix < blockSize; ++ix) {
		buffer[ix] = static_cast<uint8_t>(ix % 256);
	}

	vector<thread> threads;
	threads.reserve(threadCount);

	auto timerStart = chrono::high_resolution_clock::now();

	// spawn threads
	for (auto ix = 0; ix < threadCount; ++ix) {
		threads.push_back(thread([ix, filesPerThread, blockSize, buf = buffer.get(), blocksPerFile, &baseDirectory]() {
			worker(baseDirectory, ix, filesPerThread, blockSize, buf, blocksPerFile);
		}));
	}
	for (auto&& t : threads) {
		t.join();
	}

	auto timerEnd = chrono::high_resolution_clock::now();

	auto ellapsed = timerEnd - timerStart;

	cout << "done in " << chrono::duration_cast<chrono::milliseconds>(ellapsed).count() << " ms" << endl;
	return 0;
}

class safe_descriptor{
private:
	int _innerFD;
public:
	explicit safe_descriptor(int fd) : _innerFD(fd){}

	bool valid()const{
		return _innerFD != -1;
	}
	int get()const{
		return _innerFD;
	}

	~safe_descriptor(){
		if(valid()){
			close(_innerFD);
		}
	}
};

void worker(const string & baseDir, int threadIndex, int fileCount, int blockSize, uint8_t * buf, int blocksPerFile)
{
	string threadDir = baseDir + "/thread" + to_string(threadIndex);
	// create dir
	auto status = mkdir(threadDir.c_str(), S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);
	if (status != 0 && status != EEXIST) {
		throw invalid_argument("cannot create a directory there");
	}
	// write files
	for (auto ix = 0; ix < fileCount; ++ix) {
		string file = threadDir + "/" + to_string(ix);
		auto fd = safe_descriptor(open(file.c_str(), O_WRONLY | O_CREAT | O_TRUNC , S_IRWXU | S_IRWXG | S_IRWXO));
		if (!fd.valid()) {
			throw invalid_argument("cannot create a file there");
		}
		for (auto blockIx = 0; blockIx < blocksPerFile; ++blockIx) {
			int written = 0;
			while(written < blockSize){
				written += write(fd.get(), buf+written, blockSize-written);
			}
		}
	}
	// read files
	auto readBuf = unique_ptr<uint8_t[]>(new uint8_t[blockSize]);
	for (auto ix = 0; ix < fileCount; ++ix) {
		string file = threadDir + "/" + to_string(ix);
		auto fd = safe_descriptor(open(file.c_str(), O_RDONLY));
		if (!fd.valid()) {
			throw invalid_argument("cannot open a file there");
		}
		for (auto blockIx = 0; blockIx < blocksPerFile; ++blockIx) {
			auto readSize = read(fd.get(), readBuf.get(), blockSize);
			if(readSize != blockSize || memcmp(buf, readBuf.get(), blockSize) != 0){
				cerr << "file seem corrupted";
			}
		}
	}
}
