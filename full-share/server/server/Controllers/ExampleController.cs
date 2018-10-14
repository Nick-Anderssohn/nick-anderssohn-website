using System.Collections.Generic;
using System.Linq;
using Microsoft.AspNetCore.Mvc;
using server.Models;

namespace server.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class ExampleController : ControllerBase
    {
        private readonly ExampleContext _context;

        public ExampleController(ExampleContext context)
        {
            _context = context;

            if (!_context.Examples.Any())
            {
                // Create an example if collection is empty
                _context.Examples.Add(new ExampleModel() {Value = "Chicken Dinner"});
                _context.SaveChanges();
            }
        }

        [HttpGet]
        public ActionResult<List<ExampleModel>> GetAll()
        {
            return _context.Examples.ToList();
        }

        [HttpGet("{id}", Name = "GetExample")]
        public ActionResult<ExampleModel> GetById(long id)
        {
            ExampleModel model = _context.Examples.Find(id);
            if (model == null)
            {
                return NotFound();
            }
            return model;
        }
    }
}