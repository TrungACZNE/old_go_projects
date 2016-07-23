#!/usr/bin/env python

__author__ = 'trung'

import StringIO
import json
import os
import os.path
import subprocess
import sys

class_pool = {}
primitives = [int, float, bool, str, unicode]


def fix_name(name):
    return name.title().replace("_", "")


def list_of(t):
    # t is a type
    return "[]" + t


class_name_counter = 0
def gen_generic_class_name():
    global class_name_counter
    class_name_counter += 1
    return "_Generic_Class_" + str(class_name_counter)


def gen_class(key_name=None):
    if key_name is None:
        # assumes this always return a unique name
        return gen_generic_class_name()
    else:
        # TODO: checks collision
        return fix_name(str(key_name))


def has_existed(description):
    global class_pool
    for k, v in class_pool.iteritems():
        if v == description:
            return k
    return None


def class_of(description, key_name=None):
    global class_pool
    existing = has_existed(description)
    if existing is not None:
        return existing
    else:
        new_class = gen_class(key_name)
        class_pool[new_class] = description
        return new_class


def unpack_primitive(obj):
    if type(obj) == unicode:
        return str.__name__
    return type(obj).__name__


def unpack_list(obj, key_name=None):
    if len(obj) > 0:
        elem = obj[0]
        return list_of(unpack(elem, key_name))
    else:
        return None


def unpack_dict(obj, key_name=None):
    # Generate a new class describing obj
    # returns the class' name
    sorted_keys = sorted(obj.keys())
    description = {}
    for key in sorted_keys:
        t = unpack(obj[key], key)
        if t is not None:
            description[key] = t
    return class_of(description, key_name)


def unpack(obj, key_name=None):
    if type(obj) in primitives:
        return unpack_primitive(obj)

    elif type(obj) == list:
        # assumes all elements in this list are of the same type
        return unpack_list(obj, key_name)

    elif type(obj) == dict:
        return unpack_dict(obj, key_name)

    else:
        return None


GO_CLASS_TEMPLATE = """type %s struct {
%s
}"""
GO_TYPE_TEMPLATE = r'%s %s `json:"%s"`'


def key_name_to_go_field(key_name):
    return key_name.title().replace("_","")


def to_go():
    global class_pool
    result = ""
    l = ["Root"] + [x for x in sorted(class_pool.keys()) if x != "Root"]
    for class_name in l:
        class_description = class_pool[class_name]
        class_body = ""
        for key in sorted(class_description.keys()):
            key_type = class_description[key]
            go_type = ""
            if key_type == "int":
                go_type = "int64"
            elif key_type == "float":
                go_type = "float64"
            elif key_type == "bool":
                go_type = "bool"
            elif key_type == "str":
                go_type = "string"
            else: # class
                go_type = key_type
            class_body += "\t" + GO_TYPE_TEMPLATE % (
                key_name_to_go_field(key),
                go_type,
                key
            ) + "\n"
        result += GO_CLASS_TEMPLATE % (
            class_name,
            class_body
        ) + "\n"
    return result


def flatten_dict(obj, key_name):
    sorted_keys = sorted(obj.keys())
    description = {}
    for key in sorted_keys:
        t = flatten(obj[key], key)
        if t is not None:
            if type(t) == dict:
                for childkey, val in t.iteritems():
                    new_name = key + "_" + childkey
                    description[new_name] = val
            else:
                description[key] = t
    return description


def flatten(obj, key_name):
    if type(obj) in primitives:
        val = unpack_primitive(obj)

    elif type(obj) == list:
        print "Can't flatten lists"
        sys.exit(1)

    elif type(obj) == dict:
        val = flatten_dict(obj, key_name)

    else:
        val = None
    return val


def main():
    blob = sys.stdin.read()
    obj = json.loads(blob)
    # unpack(obj, "root")

    if len(sys.argv) > 1 and sys.argv[1] == "--flatten" :
        class_pool["Root"] = flatten(obj, "")
    else :
        unpack(obj, "root")

    output_blob = to_go()

    if "GOROOT" in os.environ:
        gofmt = os.environ["GOROOT"] + "/bin/gofmt"
        if os.path.isfile(gofmt):
            p = subprocess.Popen([gofmt], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            stdout, stderr = p.communicate(output_blob)
            print stdout
    else:
        print output_blob

if __name__ == '__main__':
    main()
