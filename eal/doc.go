/*
Package eal wraps EAL initialization and provides some additional functionality
on top of that. Every CPU's logical core which is setup by EAL runs its own
function which essentially receives functions to execute via Go channel. So you
may run arbitrary Go code in the context of EAL thread.

EAL may be initialized via command line string, parsed command line string or a
set of Options.

Please note that some functions may be called only in EAL thread because of TLS
(Thread Local Storage) dependency.

API is a subject to change. Be aware.
*/
package eal
