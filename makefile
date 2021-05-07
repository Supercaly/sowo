GO 			:= go
BIN     	:= ./bin
EXAMPLE     := ./examples
EXAMPLES    := $(wildcard $(EXAMPLE)/*.sowo)
BINS    	:= $(patsubst $(EXAMPLE)/%.sowo,$(BIN)/%.c,$(EXAMPLES))
EXES    	:= $(patsubst $(BIN)/%.c,$(BIN)/%,$(BINS))
TOKENS    	:= $(patsubst $(BIN)/%.c,$(BIN)/%_tok.txt,$(BINS))
ASTS    	:= $(patsubst $(BIN)/%.c,$(BIN)/%_ast.json,$(BINS))

.PHONY: all

all: examples exec_examples

examples: $(BINS) | $(BIN) $(EXES)

$(BIN)/%.c: $(EXAMPLE)/%.sowo | $(BIN)
	$(GO) run . $< -o $@ --save-tokens --save-ast

$(BIN):
	mkdir $@

$(EXES): $(BINS) | $(BIN)
	$(CC) -Wall $< -o $@

exec_examples: $(EXES)
	$<

clean_examples:
	rm -f $(BINS) $(EXES) $(ASTS) $(TOKENS)
