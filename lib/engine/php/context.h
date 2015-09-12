#ifndef CONTEXT_H
#define CONTEXT_H

#define BUFFER_SIZE 16 // Buffer size to use for the engine output buffer, in bytes.

typedef struct _engine_context {
	php_engine *engine; // Parent engine instance.
	struct {
		char *bytes;
		size_t size;
		size_t used;
	} buffer;
} engine_context;

engine_context *context_new(php_engine *engine);
void context_run(engine_context *context, char *filename);
size_t context_sync(engine_context *context, char **buffer);
void context_destroy(engine_context *context);

#endif