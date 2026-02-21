from typing import TypedDict, NotRequired
from itertools import permutations
from probability import distribute_balls_unlabeled, unique_permutation_count

class Item(TypedDict):
  name:str
  quantity:int
  unstackable:NotRequired[bool]

# encodes an item into a string for easy comparisons
def encode_item(item:Item):
  return f"{item['name']}//{item['quantity']}//{item.get("unstackable", False)}"

# given a full encoded inventory, figure out the different arrangements and return them as encoded items
def calculate_inventory_arrangements(inventory:list[str]):
  # return set(permutations(inventory))
  return unique_permutation_count(inventory)

# given a stack size, and number of slots to spread it to, return all the possible combinations (order ignored)
def spread_stack(stack:int, slots:int):
  return distribute_balls_unlabeled(stack, slots)

# TODO handle multiple stacks of different items
def calculate_inventories(items:list[Item]):
  de_blanked = [item for item in items if item["name"] != "Blank"]
  item_names = [item["name"] for item in de_blanked]
  if len(de_blanked) != len(set(item_names)):
    raise Exception("No duplicate items allowed")
  inventories = {",".join([encode_item(item) for item in items])}
  
  blanks = len(items) - len(de_blanked)
  for item in de_blanked:
    if(item["quantity"]<2):
      continue
    for i in range(blanks):
      slots = i+2
      stack = item["quantity"]
      spreads = spread_stack(stack, slots)
      for spread in spreads:
        inventories.add(",".join(
          [encode_item(spread_item) for spread_item in de_blanked if (spread_item["name"] != item["name"])] +
          [f"{item['name']}//{spread_count}" for spread_count in spread]+
          ["Blank//1//True" for _ in range(blanks-i-1)]
        ))
  return inventories

  

items:list[Item] = [
  {
    "name":"Axe",
    "quantity":1
  },
  {
    "name":"Hoe",
    "quantity":1
  },
  {
    "name":"WateringCan",
    "quantity":1
  },
  {
    "name":"Pickaxe",
    "quantity":1
  },
  {
    "name":"MeleeWeapon",
    "quantity":1
  },
  {
    "name":"Parsnip Seeds",
    "quantity":15
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  }
]

count = 0
inventories = calculate_inventories(items)
abstract_inventories = {}
for ii,i in enumerate(inventories):
  print(f"Checking Inventory {ii+1}/{len(inventories)}")
  inventory = i.split(",")
  abstract_inventory_dict = {}
  for item in inventory:
    if item not in abstract_inventory_dict.keys():
      abstract_inventory_dict[item] = 0
    abstract_inventory_dict[item] += 1
  abstract_inventory = ",".join(str(v) for v in [sorted(abstract_inventory_dict.values())])
  if abstract_inventory in abstract_inventories.keys():
    i_count = abstract_inventories[abstract_inventory]
  else:
    # i_count = len(calculate_inventory_arrangements(inventory))
    i_count = calculate_inventory_arrangements(inventory)
    abstract_inventories[abstract_inventory] = i_count
  print(f"Found {i_count} permutations")
  count += i_count

print(f"Total: {count}")