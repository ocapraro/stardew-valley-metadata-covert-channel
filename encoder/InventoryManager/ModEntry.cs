using StardewModdingAPI;
using StardewModdingAPI.Events;
using StardewValley;
using StardewValley.Menus;
using StardewValley.Objects;

namespace InventoryManager;

public class ModEntry : Mod
{
  public override void Entry(IModHelper helper)
  {
    Monitor.Log("Mod loaded!", LogLevel.Info);
    helper.Events.Display.MenuChanged += OnMenuChanged;
  }

  /// <summary>
  /// Determines whether the current menu transition indicates the player is attempting to sleep.
  /// </summary>
  /// <param name="e">Menu change event arguments.</param>
  /// <returns>
  /// <c>true</c> when the world is ready, the new menu is a dialogue in the farmhouse,
  /// and the player is positioned at the expected bed interaction tiles; otherwise <c>false</c>.
  /// </returns>
  private static bool IsSleeping( MenuChangedEventArgs e) {
    if (!Context.IsWorldReady)
          return false;

      if (e.NewMenu is null)
          return false;

      if (e.NewMenu is not DialogueBox)
        return false;

      if (Game1.currentLocation?.Name != "FarmHouse")
        return false;

      var bed = Game1.currentLocation.furniture.FirstOrDefault(f => f is BedFurniture);

      if (bed is null)
        return false;
      var playerLocation = Game1.player.TilePoint.ToVector2();
      var bedLocation = bed.TileLocation;
      
      if 
      (
        playerLocation.Y != bedLocation.Y+1 || 
        (
          playerLocation.X != bedLocation.X &&
          playerLocation.X != bedLocation.X+1
        )
      )
        return false;
      return true;  
  }

  private void OnMenuChanged(object? sender, MenuChangedEventArgs e)
  {
    if(IsSleeping(e))
    {
      var playerInventory = Game1.player.Items;
      
      Monitor.Log($"Inventory slots: {playerInventory.Count}", LogLevel.Info);
      for (var slot = 0; slot < playerInventory.Count; slot++)
      {
        var item = playerInventory[slot];
        if (item is null)
          continue;

        Monitor.Log($"Slot {slot}: {item.QualifiedItemId}, {item.DisplayName}, {item.Stack}", LogLevel.Info);
      }
    }    
  }
}