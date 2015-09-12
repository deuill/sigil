// Standard library.
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

// PHP includes.
#include <main/php.h>
#include <main/SAPI.h>
#include <main/php_main.h>

// Local includes.
#include "engine.h"
#include "context.h"

engine_context *context_new(php_engine *engine, void *parent) {
	engine_context *context;

	#ifdef ZTS
		void ***tsrm_ls = engine->tsrm_ls;
	#endif

	// Initialize context.
	context = (engine_context *) malloc((sizeof(engine_context)));
	if (context == NULL) {
		return_multi(NULL, 1);
	}

	context->engine = engine;
	context->parent = parent;

	SG(server_context) = (void *) context;

	// Initialize request lifecycle.
	if (php_request_startup(TSRMLS_C) == FAILURE) {
		SG(server_context) = NULL;
		free(context);

		return_multi(NULL, 1);
	}

	return_multi(context, 0);
}

void context_run(engine_context *context, char *filename) {
	#ifdef ZTS
		void ***tsrm_ls = context->engine->tsrm_ls;
	#endif

	// Attempt to execute script file.
	zend_first_try {
		zend_file_handle script;

		script.type = ZEND_HANDLE_FILENAME;
		script.filename = filename;
		script.opened_path = NULL;
		script.free_filename = 0;

		php_execute_script(&script TSRMLS_CC);
	} zend_end_try();

	return_multi(NULL, 0);
}

void context_destroy(engine_context *context) {
	php_request_shutdown((void *) 0);

	SG(server_context) = NULL;
	free(context);
}