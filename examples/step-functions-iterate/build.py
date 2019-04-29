#!/usr/bin/python
CLOUDFORMATION_OUTPUT_FILE = "output.yaml"
CLOUDFORMATION_INPUT_FILE = "step_functions_iterate.yaml"


def replace(old_string, new_content_file_name, input_file_name, output_file_name):
    with open(input_file_name) as input_file:
        with open(output_file_name, "w") as output_file:
            with open(new_content_file_name, "r") as new_content_file:
                replace_maintain_indentation(
                    old_string,
                    new_content_file,
                    input_file,
                    output_file)


def replace_maintain_indentation(old_string, new_content_file, input_file, out_file):
    datafile = input_file.readlines()
    for line in datafile:
        if old_string in line:
            indentation = line.find(old_string)
            copy_with_indentation(new_content_file, out_file, indentation)
        else:
            out_file.write(line)


def copy_with_indentation(in_file, out_file, indentation):
    datafile = in_file.readlines()
    for line in datafile:
        if "\n" in line:
            out_file.write(' ' * indentation + line)
        else:
            out_file.write(' ' * indentation + line + "\n")


replace("{{stateMachine}}", "state-machine.json", CLOUDFORMATION_INPUT_FILE, CLOUDFORMATION_OUTPUT_FILE)
