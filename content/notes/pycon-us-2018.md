Title: PyCon US 2018 - Main conference
Date: 2018-05-11 12:00
Category: software

# Keynote: Dan Callahan

## Notes

- Platform detected the tool
- Rust -> WASM -> Web
- Rust -> DLL -> Python
- The Birth and Death of Javascript, Gary Berhardt


# Thinking Outside the GIL with AsyncIO and Multiprocessing

[Talk overview](https://us.pycon.org/2018/schedule/presentation/103/)


## Notes

Previously:

- Thread pool for I/O
- Underutilized CPU
- Memory usage

Py3 switch:

- ~45% memory savings
- 20% runtime reduction

Multiprocessing

- CPU scaling
- CPU utilization
- Automatic IPC
- Pool.map is really useful
- One task per processor
- Memory duplication (gc)

AsyncIO

- Based on futures
- Faster than threads
- Massive I/O concurrency
- Still relyin on one process
- Process still limited by GIL

Ideally: Asyncio + multiprocessing

- Event loop per process
- Queues for work/results
- asyncio.wait

Considerations

- Minimize what you pickle
- Prechunk work items
- Aggregate results in the child
- Use map/reduce

Projects

- [aiomultiprocess](https://github.com/jreese/aiomultiprocess)

Great tools make complex tasks simple




# The AST and Me

[Talk overview](https://us.pycon.org/2018/schedule/presentation/107/)

Friday 3:15 p.m.–4 p.m. in Grand Ballroom B

## Notes


# Trio: Async concurrency for mere mortals

[Talk overview](https://us.pycon.org/2018/schedule/presentation/163/)

# Love your bugs

[Talk overview](https://us.pycon.org/2018/schedule/presentation/156/)

- Bitflips are real
- Server side gating and feature flags
- Debugging requires a growth mindset


# Elegant Solutions For Everyday Python Problems

Friday 5:10 p.m.–5:40 p.m. in Room 26A/B/C

[Talk overview](https://us.pycon.org/2018/schedule/presentation/164/)

## Notes


# Inside the Cheeseshop: How Python Packaging Works

[Talk overview](https://us.pycon.org/2018/schedule/presentation/148/)

[Video](https://www.youtube.com/watch?v=AQsZsgJ30AE)


# Dataclasses: The code generator to end all code generators

[Talk overview](https://us.pycon.org/2018/schedule/presentation/94/)

Saturday 3:15 p.m.–4 p.m. in Grand Ballroom C


https://us.pycon.org/2018/schedule/presentation/80/

