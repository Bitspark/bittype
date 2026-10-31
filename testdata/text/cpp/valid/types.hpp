#include <vector>
namespace api {
struct Person { int id; };
template<typename T, class U = T>
struct Page { std::vector<T> items; U extra; };
using People = Page<Person>;
enum class Mode : int { Read = 1, Write = 2, Minus = -1 };
enum Plain : int;
struct Base { virtual ~Base() noexcept = 0; };
template<class T> class Reader : public virtual Base {
public:
  virtual T read(const char* key) const & noexcept = 0;
};
struct Shapes {
  int *array[3];
  int (*pointer)[3];
  int (*callback)(double) noexcept;
  int Person::* member;
  const int* pointerToConst;
  int* const constPointer;
  int&& move;
  int& reference;
  volatile int flag;
  using Sig = void(int, ...);
};
using Nested = std::vector<std::vector<int>>;
template<typename T> using Value = typename T::value_type;
template<typename T> using NestedType = typename T::template Rebind<int>::type;
struct Forward;
}
