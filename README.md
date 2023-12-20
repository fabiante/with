# with

This package aims to simplify handlign of several interface types like `io.Closer`.

The idea is essentially to provide something like the `with` keyword in Python:

```python
with open('file.txt', 'w') as file:
    file.write('Hello Python')
```

This package is currently only an exploration of mine. Use it if you finde it useful.
I am not yet sure if this library has a real value or is simply something that one
would manually do if needed.
