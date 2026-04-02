using StardewModdingAPI;
using StardewModdingAPI.Events;
using StardewValley;
using StardewValley.Menus;
using StardewValley.Objects;
using System.Net.Http.Json;
using System.Text.Json;

namespace InventoryManager;


public class ResponseItem
{
  public string? Name { get; set; }
  public int Stack { get; set; }
}

public class Response
{
  public List<ResponseItem>? Items { get; set; }
}
public class ModEntry : Mod
{
  private static readonly HttpClient HttpClient = new();
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

  private async void HandleSleep()
  {
    Monitor.Log("Player sleeping", LogLevel.Info);
    var playerInventory = Game1.player.Items;

    var requestInventory = new Response
    {
      Items = playerInventory.Select(item =>
        item is null
          ? new ResponseItem { Name = "Blank", Stack = 0 }
          : new ResponseItem { Name = item.DisplayName, Stack = item.Stack }
      ).ToList()
    };

    var setInventoryResponse = await HttpClient.PostAsJsonAsync("http://localhost:8080/setInventory", requestInventory);
    if (!setInventoryResponse.IsSuccessStatusCode)
    {
      Monitor.Log($"Failed to set inventory: {(int)setInventoryResponse.StatusCode}", LogLevel.Error);
      return;
    }

    var httpResponse = await HttpClient.GetAsync("http://localhost:8080/encode");
    if (!httpResponse.IsSuccessStatusCode)
    {
      Monitor.Log($"Failed to encode inventory: {(int)httpResponse.StatusCode}", LogLevel.Error);
      return;
    }

    var json = await httpResponse.Content.ReadAsStringAsync();
    var result = JsonSerializer.Deserialize<Response>(json);
    if (result?.Items is null)
      return;
    // Record all current items
    var originalItems = new List<Item>(playerInventory);
    // Clear the hotbar
    playerInventory.Clear();
    
    // For each item in the response
    for (var i = 0; i < result.Items.Count; i++)
    {
      var responseItem = result.Items[i];
      
      // Add blank slots as null
      if (responseItem.Name == "Blank")
      {
        playerInventory.Add(null);
        continue;
      }
      
      // Find the original item by name
      var originalItem = originalItems.FirstOrDefault(item => item is not null && item.DisplayName == responseItem.Name);
      
      if (originalItem is not null)
      {
        // Copy the item
        var copiedItem = originalItem.getOne();
        // Set the stack
        copiedItem.Stack = responseItem.Stack;
        // Add to hotbar
        playerInventory.Add(copiedItem);
      }
      else
      {
        playerInventory.Add(null);
      }
    }
    // Monitor.Log(res., LogLevel.Info);
  }

  private void OnMenuChanged(object? sender, MenuChangedEventArgs e)
  {
    Monitor.Log("Menu Opened", LogLevel.Info);
    if(IsSleeping(e))
    {
      var playerInventory = Game1.player.Items;
      HandleSleep();
      // Monitor.Log($"Inventory slots: {playerInventory.Count}", LogLevel.Info);
      // for (var slot = 0; slot < playerInventory.Count; slot++)
      // {
      //   var item = playerInventory[slot];
      //   if (item is null)
      //     continue;

      //   Monitor.Log($"Slot {slot}: {item.QualifiedItemId}, {item.DisplayName}, {item.Stack}", LogLevel.Info);
      // }
    }    
  }
}