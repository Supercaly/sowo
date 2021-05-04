GO 			:= go
BIN     	:= ./bin
EXAMPLE     := ./examples
EXAMPLES    := $(wildcard $(EXAMPLE)/*.sowo)
BINS    	:= $(patsubst $(EXAMPLE)/%.sowo,$(BIN)/%.c,$(EXAMPLES))
EXES    	:= $(patsubst $(BIN)/%.c,$(BIN)/%,$(BINS))

.PHONY: all

all: examples

examples: $(BINS) | $(BIN) $(EXES)

$(BIN)/%.c: $(EXAMPLE)/%.sowo | $(BIN)
	$(GO) run . $< -o $@

$(BIN):
	mkdir $@

$(EXES): $(BINS) | $(BIN)
	$(CC) -Wall $< -o $@

clean_examples:
	rm $(BINS) $(EXES)

